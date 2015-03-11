package go;

public class IntWrapper implements Ptr_int{

	Int addr;

	public IntWrapper(Int value){
		this.addr = value;
	}

	public boolean equals(Object o){
		return o instanceof IntWrapper && ((IntWrapper)o).addr == this.addr;
	}

	public int value(){return this.addr.value;}

	public void set(int v){this.addr.set(v);}

}
