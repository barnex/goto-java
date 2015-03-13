
public class Int implements Ref_int{

	public int value;

	public Int(){
	}

	public Int(int value){
		this.value = value;
	}

	public int value(){
		return value;
	}

	public void set(int value){
		this.value = value;
	}

	public boolean equals(Object other){
		return other instanceof Int 
			&& ((Int)other).value == this.value;	
	}

	public void inc(){
		this.value++;
	}
}
