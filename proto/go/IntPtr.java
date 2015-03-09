package go;

public abstract class IntPtr{
	
	public abstract void set(int v);
	public abstract int value();
	public abstract boolean equals(Object other);
	public void inc(){ set(value()+1); }
	public void dec(){ set(value()-1); }

	public void add(int x){ set(value() + x); }
	public void sub(int x){ set(value() - x); }
	public void mul(int x){ set(value() * x); }
	public void quo(int x){ set(value() / x); }
	public void rem(int x){ set(value() % x); }
	public void and(int x){ set(value() & x); }
	public void or (int x){ set(value() | x); }
	public void xor(int x){ set(value() ^ x); }
	public void shl(int x){ set(value()<< x); }
	public void shr(int x){ set(value()>> x); }
	public void andnot(int x){ set(value() & ~x) ;}
}
