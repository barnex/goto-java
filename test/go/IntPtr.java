package go;

public abstract class IntPtr{
	
	public abstract void set(int v);
	public abstract int value();
	public void inc(){ set(value()+1); }
	public void dec(){ set(value()-1); }

}
